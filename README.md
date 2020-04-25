# goldmark-wikilinks

This is a goldmark plugin for processing wikilinks
(links that look like \[\[Hi there]]). You can add it
directly to your goldmark to have it turn links like
the example into `[Hi there](Hi there.html)`. Using a `FilenameNormalizer` you can
customize what the link destination is. Using a `WikilinkTracker` you can gather up
the links to create a list of backlinks, which was the reason I created this in
the first place.