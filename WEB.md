# ekzyis.com

My personal website including blog.

## Development

This site consists of only static HTML, CSS, JS in public/.

The files are built (or "rendered") with the golang [`text/template`](https://pkg.go.dev/text/template) standard package. It doesn't use [`html/template`](https://pkg.go.dev/html/template) since I had problems including HTML like a common header, navigation menu, footer for a reusable layout. But this shouldn't be a problem since there is no user-generated content (yet?).

To build the files, a [Makefile](./Makefile) is used.

Run `make build` to create the `renderer` binary.

Run `make render` to render all files in public/.

Deployment is done by rendering all files in production mode and then copying them where a webserver like `nginx` can serve them.

I use [`deploy.sh`](./deploy.sh) for this.

## How to create new blog post

1. Create new Markdown file in blog/
2. It needs to have this header:

```
Title:        title
Date:         date
ReadingTime:  time
Sats:         0
Comments:     comments

---
```

3. Update `ComputeTitle` in `html.go` (TODO: make this no longer required)
4. Run `make render`.

Done!
