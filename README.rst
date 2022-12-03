===============
wordleguesses
===============
--------------------------------------------------------------------------------
A command-line tool to easily generate a list of candidate Wordle guesses
--------------------------------------------------------------------------------

Introduction
=============

When playing Wordle, sometimes it can be helpful to write out a list
of candidate guesses.  For example, you might be considering all the
possibilities that arise from changing the first character in the
sequence "?A_AM". Assuming through previous play you've already ruled
out 'R', 'I', 'S', 'E', 'N', 'G', 'Y', 'C', 'U', and 'K', you would
end up generating a list like (assuming you are being exhaustive and
not skipping improbable possibilities)::

  AA_AM
  BA_AM
  DA_AM
  FA_AM
  HA_AM
  JA_AM
  LA_AM
  MA_AM
  OA_AM
  PA_AM
  QA_AM
  TA_AM
  VA_AM
  WA_AM
  XA_AM
  ZA_AM

Writing lists like these out can be quite laborious, and can create a
significant hindrance to those with diminished
dexterity. **wordle_guesses** is intended to alleviate this burden by
listing out those candidate guesses for you.

Usage
=====

The basic usage of **wordle_guesses** is::

  go run main.go <template>

where *template* is composed of 5 characters that are letters
(characters from the alphabet like 'A' and 't'), the character '_'
(underscore) and the character '.' (period). ('.' is used instead of
'?' to avoid issues with command-line processors that try to perform
substitution using '?'.)

So, an example close to that in the Introduction would be::

  go run main.go .A_AM

Since typing in lots of uppercase characters is a pain,
**wordle_guesses** ignores the case of the characters in the supplied
template, so you can instead use::

  go run main.go .a_am

You can specify the set of letters to ignore using the
**-e** (exclude) option following by a list of the letters to ignore. Using
the example from the introduction, you could have something like
this::

  go run main.go -e risengycuk .a_am

There are few more options that you can use. For the complete usage,
use **wordle_guesses** with the  **-d** (description) option.
