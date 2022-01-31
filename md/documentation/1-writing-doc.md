# Writing documentation

## Titles and nav

The sidebar navigation is generated from the structure of the root `md`-directory.
Titles are read from one of two places.

For files, the title is read from the first line of the file, prefixed with a single `#`.
For example for this file the title read is `Writing documentation`.

For directories, each directory contains a file called `_title`.
This file contains the title in plain text.
For example for this directory the `_title`-file reads `Documentation`.

For the project title (the one in the top left), similar to normal directories, there is a `_title`-file in the root `md`-directory.

All directories and files will automatically get an entry in the sidebar navigation.

## Usage

This project sets up the a documentation server on the GIN-router in `main.go`. This is done by a few lines of code in the project's main-function:
```go
// create doc-handler
docSrv := doc.MustCreateHandler(docfs.FS(false))
// register doc-routes
docSrv.MustRegisterGin(r)
```

`docfs.FS` is a virtual filesystem generated from the `md` folder in the project's root folder.

Make-goals to use while writing documentation:
* [`watch-doc`](#/building/makefile#watchdoc): starts `main.go` and restarts if any `md`- or `_title`-file in the `md` directory is updated
* [`generate-doc`](#/building/makefile#generatedoc): generates `internal/docfs/resources.go` containing all files in the `md` directory

## Links

Links in markdown are created using the following syntax:
```markdown
[linktext](link)
```

The `link` can be a link to an external site, a local page and/or a hashlink to a heading on a page.

### External links
External links are created like [this](https://trondheim.kommune.no/).
```markdown
[this](https://trondheim.kommune.no/)
```

### Page links
Links to other pages are created similarly to external links, except the host-location can be omitted, only including the url-hash.

[this](#/building/makefile) is an example:
```markdown
[this](#/building/makefile)
```

### Hashlinks
A hashlink can be appended to an external site, a local page, or be on it's own. They look like [this](#page-links):
```markdown
[this](#page-links)
```

They can also be appended to pagelinks like [so](#/building/makefile#live-reloading-watch):
```markdown
[so](#/building/makefile#live-reloading-watch)
```

A hashlink links to a heading `id`.
The IDs for this documentation page are automatically generated in the following way:
```javascript
const headContent = "Live reloading (watch)";
const hashLinkID = headContent
  // replace all uppercase characters with their corresponding lower case variant
  .toLowerCase()
  // string is now "live reloading (watch)"
  // remove everything that is not an alphanumeric or a space
  .replace(/[^a-z0-9 ]+/g, '')
  // string is now "live reloading watch"
  // replace spaces with dashes
  .replace(/ /g, '-');
  // resulting string is "live-reloading-watch"
```

NB: Keep in mind that since all non alphanumeric characters are removed, dashes are also removed.
So the heading `Watch-all` would be renamed to `watchall`.
