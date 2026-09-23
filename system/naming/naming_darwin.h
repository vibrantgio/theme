// The C surface of the package's Objective-C (naming_darwin.m). The Go side
// (naming.go) states what the order is and where it is answered without the
// platform.

#ifndef VGIO_THEME_SYSTEM_NAMING_NAMING_DARWIN_H
#define VGIO_THEME_SYSTEM_NAMING_NAMING_DARWIN_H

// vgio_naming_compare returns -1, 0 or 1 for a ordered before, alike with or
// after b, by NSString's localizedStandardCompare: — the comparison the
// platform's file browser sorts by, which reads a run of digits as the
// number it spells and forces an order between names that collate alike.
//
// Both arguments are UTF-8. A byte sequence that is not UTF-8 is compared as
// the empty name rather than refused, so a listing with one unreadable name
// still sorts.
//
// No window and no NSApplication are required, and none is created: the
// comparison touches no application state, so any thread may call it.
extern int vgio_naming_compare(const char *a, const char *b);

#endif
