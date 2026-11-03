# Vulnetix Terraform provider
#
# Releases are signed with a GPG key whose public half is registered on
# registry.terraform.io and the OpenTofu registry. The key lives in the local
# keyring; its fingerprint and passphrase live in .env.production (gitignored).
#
#   just gpg-key      create the signing key (once) and record it in .env.production
#   just secrets      push GPG_PRIVATE_KEY and PASSPHRASE to the repo's Actions secrets
#   just public-key   print the armoured public key to paste into the registries
#   just release 1.0.0  tag v1.0.0 and push it, which runs the Release workflow

repo := "Vulnetix/terraform-provider-vulnetix"
key_uid := "Vulnetix Terraform Provider <chris@vulnetix.com>"

default:
    @just --list

# Create the RSA 4096 signing key non-interactively, unless .env.production already names one in the keyring
gpg-key:
    #!/usr/bin/env bash
    set -euo pipefail
    touch .env.production && chmod 600 .env.production
    source .env.production
    if [[ -n "${GPG_FINGERPRINT:-}" ]] && gpg --batch --list-secret-keys "$GPG_FINGERPRINT" >/dev/null 2>&1; then
        echo "✅ Signing key already exists: $GPG_FINGERPRINT"
        exit 0
    fi
    PASSPHRASE="${PASSPHRASE:-$(openssl rand -hex 32)}"
    # The registries accept RSA or DSA only, not gpg's default ECC.
    gpg --batch --pinentry-mode loopback --passphrase-fd 0 \
        --quick-gen-key "{{key_uid}}" rsa4096 sign never <<<"$PASSPHRASE"
    fpr=$(gpg --batch --with-colons --list-secret-keys "{{key_uid}}" \
        | awk -F: '/^sec:/{s=1;next} s&&/^fpr:/{f=$10;s=0} END{print f}')
    sed -i '/^\(PASSPHRASE\|GPG_FINGERPRINT\)=/d' .env.production
    printf 'GPG_FINGERPRINT=%s\nPASSPHRASE=%s\n' "$fpr" "$PASSPHRASE" >> .env.production
    echo "✅ Created signing key $fpr (recorded in .env.production)"

# Set the GPG_PRIVATE_KEY and PASSPHRASE Actions secrets from the local key
secrets:
    #!/usr/bin/env bash
    set -euo pipefail
    source .env.production
    : "${GPG_FINGERPRINT:?not set in .env.production, run just gpg-key}"
    : "${PASSPHRASE:?not set in .env.production, run just gpg-key}"
    gpg --batch --pinentry-mode loopback --passphrase-fd 0 \
        --armor --export-secret-keys "$GPG_FINGERPRINT" <<<"$PASSPHRASE" \
        | gh secret set GPG_PRIVATE_KEY -R {{repo}}
    printf '%s' "$PASSPHRASE" | gh secret set PASSPHRASE -R {{repo}}
    gh secret list -R {{repo}}

# Print the armoured public key for registry.terraform.io and the OpenTofu registry
public-key:
    #!/usr/bin/env bash
    set -euo pipefail
    source .env.production
    gpg --batch --armor --export "${GPG_FINGERPRINT:?run just gpg-key}"

# Tag vVERSION on origin/main and push it, which runs the Release workflow
release version:
    #!/usr/bin/env bash
    set -euo pipefail
    [[ "{{version}}" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "❌ version must be X.Y.Z"; exit 1; }
    git fetch -q origin
    git tag -a "v{{version}}" -m "v{{version}}" origin/main
    git push origin "v{{version}}"
