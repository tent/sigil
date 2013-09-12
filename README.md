# Sigil

![Sigil's Sigil](https://sigil.cupcake.io/Sigil)
![Sigil's inverted Sigil](https://sigil.cupcake.io/Sigil?inverted=1)

Sigil is a deterministic [identicon](https://en.wikipedia.org/wiki/Identicon) generator originally created for [Cupcake](https://cupcake.io).

Sigil creates a recognizable* identicon for each user that can be consistently used across many sites when a user-uploaded avatar or Gravatar is unavailable or unwanted. Sigils can be generated by any app or service using this application or other implementation of the Sigil protocol. Sites can also use Sigil [generation services](https://sigil.cupcake.io). 

\* see _Privacy_ below

### Using Sigil

Sigil 0.1 creates identicons to represent _users_. A future release will have additional algorithms and representations for non-persons. This will help anyone who encounters Sigils to determine quickly if they represent users or other content.

Sigil generates an identicon from a source string. In most cases(see _Privacy_) the source string will be the user's email address. Query parameters can be used to request an SVG or PNG of any size and switch the foreground and background colors if desired. Reversing the background and foreground colors does not change the information represented in a Sigil, so the reversed image can be used if the style is more compatible with the host app.

You can generate Sigils internally using this app (or your own implementation of the protocol below) or use a [Sigil service](https://sigil.cupcake.io).

#### Privacy 

There exists a possibility of significant user information disclosure if used improperly. Email addresses should *not* be used as the source string (even if hashed first) if the user's email address is not meant to be publicly available through your site. It is possible to discover the email address from which the Sigil was derived. This is a serious leak if the user's address is not exposed elsewhere on the site. In these cases, use a *different* source string composed of information that is available such as a username or user id. 

### The Protocol

This repo contains a Sigil implementation written in Go, but Sigil can be implemented in any language.

Sigil uses the truncated MD5 hash of a string to create identicons. 

### Roadmap

Sigil is still a work in progress. The color scheme is still being finalized and this implementation may be optimized further. Additional background colors will likely also be provided. We also intend to create a separate representation for _things_ instead of _users_ i.e. favicons for Herkou apps. The suggested use instructions (source string, background color, inverted colors) may also change before the 1.0 release.

### Try it

Download this repo or visit [sigil.cupcake.io](https://sigil.cupcake.io) to generate a Sigil. 

Use Sigils in your own apps for users who don't have avatars or use your own Sigil as an avatar in other apps.