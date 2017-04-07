# clex

A 95-99% standards-compliant 
[C11](http://www.open-std.org/jtc1/sc22/wg14/www/docs/n1570.pdf) positional lexer with 
error reporting.

The lexer reports errors with (1) a message, (2) the line so far, and (3) the position of 
the occurence, so that end-user messages can be produced such as:

```
test.c:20:33: Expected 4 hexadecimal characters for universal character name
    const char* s = "Hello \u042
--------------------------------^
```

## Example Output

Executing `clex test.c test.lexemes` on:

```c
extern int puts(const char *str)

int main(int argc, char **argv) {
        puts("Hello, world!");
}
```

Produces:

```
Keyword{extern}Whitespace{ }Keyword{int}Whitespace{ }Identifier{puts}LeftParenthesis{(}
Keyword{const}Whitespace{ }Keyword{char}Whitespace{ }Star{*}Identifier{str}RightParenthesis{)}
...
```

## Path to Full Compliance

There's a couple minor outstanding issues centered upon the lexing of *punctuators*. I'm weighing the options on how to solve the issue elegantly, if possible.

The problem lexemes are:

* Ellipsis: "..."
* DoubleHash: "%:%:"

As of right now, these two punctuators will be lexed as three *Period*'s and two *Hash*'s, respectively.

The underlying issue is lack of sufficient character lookahead. Right now, it looks ahead 1-character at most. So, there's two options:

1. Extend the underlying reader to support K-character lookahead
2. Run a post-processing pass that collapses sequences of *Period* and *Hash* lexemes into *Ellipsis* and *DoubleHash*

My plan is to write up both solutions and compare them side-by-side.

## Future Plans

- Add a preprocessor
- Disambiguate integer and floating point constants into their proper sizes
- Support GCC extensions (if any modifications are necessary)
- Document the code
