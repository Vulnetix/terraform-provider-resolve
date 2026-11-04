# Every argument with a schema default, set to that default, plus the required
# arguments from config-schemas examples/LibraryItem/minimal.yaml. A plan only sets values
# a tenant left unset. See advanced.tf for the arguments that have no default.

resource "vulnetix_library_item" "pdf_tools" {
  name = "pdf-tools"
  body = "---\nname: pdf-tools\ndescription: Read, split and merge PDF files.\n---\nUse the pdf tools on the path the user gives.\n"
  type = "skill"
}
