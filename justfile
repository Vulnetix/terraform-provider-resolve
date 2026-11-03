# Vulnetix Terraform provider
#
# Releases are signed with a GPG key whose public half is registered on
# registry.terraform.io and the OpenTofu registry. The private key, its
# passphrase and fingerprint live in Bitwarden (`signing/TF_PROVIDER_GPG_*`,
# read with vx-secrets from Vulnetix/ops). No recipe needs a local keyring or a
# local file: keys are handled in a temporary GNUPGHOME that is removed on exit.
#
#   just gpg-key      create the signing key (once) and store it in Bitwarden
#   just secrets      set GPG_PRIVATE_KEY and PASSPHRASE Actions secrets from Bitwarden
#   just public-key   print the armoured public key to paste into the registries
#   just release 1.0.0  tag v1.0.0 and push it, which runs the Release workflow

repo := "Vulnetix/terraform-provider-resolve"
key_uid := "Vulnetix Terraform Provider <chris@vulnetix.com>"
ns := "signing"

default:
    @just --list

# Create the RSA 4096 signing key in a throwaway keyring and store it in Bitwarden, unless one is stored
gpg-key:
    #!/usr/bin/env bash
    set -euo pipefail
    if vx-secrets names {{ns}} | grep -qx TF_PROVIDER_GPG_PRIVATE_KEY; then
        echo "Signing key already in Bitwarden: $(vx-secrets get {{ns}} TF_PROVIDER_GPG_FINGERPRINT)"
        exit 0
    fi
    export GNUPGHOME="$(mktemp -d)"
    trap 'gpgconf --kill all 2>/dev/null || true; rm -rf "$GNUPGHOME"' EXIT
    PASSPHRASE="$(openssl rand -hex 32)"
    # The registries accept RSA or DSA only, not gpg's default ECC.
    gpg --batch --pinentry-mode loopback --passphrase-fd 0 \
        --quick-gen-key "{{key_uid}}" rsa4096 sign never <<<"$PASSPHRASE"
    fpr=$(gpg --batch --with-colons --list-secret-keys "{{key_uid}}" \
        | awk -F: '/^sec:/{s=1;next} s&&/^fpr:/{f=$10;s=0} END{print f}')
    gpg --batch --pinentry-mode loopback --passphrase-fd 0 \
        --armor --export-secret-keys "$fpr" <<<"$PASSPHRASE" \
        | vx-secrets put {{ns}} TF_PROVIDER_GPG_PRIVATE_KEY
    printf '%s' "$PASSPHRASE" | vx-secrets put {{ns}} TF_PROVIDER_GPG_PASSPHRASE
    printf '%s' "$fpr" | vx-secrets put {{ns}} TF_PROVIDER_GPG_FINGERPRINT
    echo "Created signing key $fpr (stored in Bitwarden {{ns}}/TF_PROVIDER_GPG_*)"

# Set the GPG_PRIVATE_KEY and PASSPHRASE Actions secrets straight from Bitwarden
secrets:
    vx-secrets to-gh {{ns}} TF_PROVIDER_GPG_PRIVATE_KEY --repo {{repo}} --name GPG_PRIVATE_KEY
    vx-secrets to-gh {{ns}} TF_PROVIDER_GPG_PASSPHRASE --repo {{repo}} --name PASSPHRASE
    gh secret list -R {{repo}}

# Print the armoured public key for registry.terraform.io and the OpenTofu registry
public-key:
    #!/usr/bin/env bash
    set -euo pipefail
    export GNUPGHOME="$(mktemp -d)"
    trap 'gpgconf --kill all 2>/dev/null || true; rm -rf "$GNUPGHOME"' EXIT
    fpr="$(vx-secrets get {{ns}} TF_PROVIDER_GPG_FINGERPRINT)"
    vx-secrets get {{ns}} TF_PROVIDER_GPG_PASSPHRASE \
        | gpg --batch --quiet --pinentry-mode loopback --passphrase-fd 0 \
            --import <(vx-secrets get {{ns}} TF_PROVIDER_GPG_PRIVATE_KEY) 2>/dev/null
    have=$(gpg --batch --with-colons --list-keys | awk -F: '/^pub:/{s=1;next} s&&/^fpr:/{print $10; exit}')
    [[ "$have" == "$fpr" ]] || { echo "stored key fingerprint $have does not match $fpr" >&2; exit 1; }
    gpg --batch --armor --export "$fpr"

# Tag vVERSION on origin/main and push it, which runs the Release workflow
release version:
    #!/usr/bin/env bash
    set -euo pipefail
    [[ "{{version}}" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "❌ version must be X.Y.Z"; exit 1; }
    git fetch -q origin
    git tag -a "v{{version}}" -m "v{{version}}" origin/main
    git push origin "v{{version}}"
