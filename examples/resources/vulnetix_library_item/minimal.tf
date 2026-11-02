# Only the arguments the kind requires. From config-schemas examples/LibraryItem/minimal.yaml.

resource "vulnetix_library_item" "pdf_tools" {
  name = "pdf-tools"
  body = "---\nname: pdf-tools\ndescription: Read, split and merge PDF files.\n---\nUse the pdf tools on the path the user gives.\n"
  type = "skill"
}
