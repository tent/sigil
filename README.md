# Sigil

Sigil is a deterministic [identicon](https://en.wikipedia.org/wiki/Identicon) generator. 

Sigil is intended to create a recognizable* identicon for users

We created Sigil after being inspired by Github's [recent use](https://github.com/blog/1586-identicons) of identicons, but found their code was not open source. 

\* see "Using Sigil: Privacy" below

## Using Sigil

Sigil 0.1 is designed to represent _users_. A future release will have a different algorithm and representation for non-persons. This will help anyone who encounters Sigils to determine quickly if they represent users or other content.

### Privacy 

There exists a possibility of significant user information disclosure if used improperly. Email addresses should *not* be used as the source string (even if hashed first) if the user's email address is not meant to be publicly available through your site. It is possible to discover the email address from which the Sigil was derived. This is a serious leak if the user's address is not exposed elsewhere on the site. In these cases, use a *different* source string composed of information that is available such as a username or user id. 

## The Protocol

This repo contains a Sigil implementation written in Go, but Sigil can be implemented in any language.

Sigil uses the truncated MD5 hash of a string to create identicons. 

## Roadmap

Sigil is still a work in progress. The color scheme is still being finalized and this implementation may be optimized further. We also intend to create a separate representation for _things_ instead of _users_ i.e. favicons for Herkou apps. The suggested use instructions (source string, background color, inverted colors) may also 

## Try it

Download this repo or visit [sigil.cupcake.io](https://sigil.cupcake.io) to generate a Sigil.