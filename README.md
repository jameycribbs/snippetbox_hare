SnippetBox example using Hare DBMS
====

This is a version of the SnippetBox web application featured in Alex
Edward's outstanding book, [Let's Go](https://lets-go.alexedwards.net/),
with Hare replacing MySQL as the DBMS.  This is just a demonstration,
mainly to show how you could use Hare in a web application.  I have removed
some of the extra code that exists in the original SnippetBox source in order
to focus mostly on the differences needed to use Hare instead of MySQL.
Surprisingly, I didn't need to change that much code to get Hare working.

[Hare](https://www.github.com/jameycribbs/hare) is a pure Go database
management system that stores each table as a text file of line-delimited
JSON.  Each line of JSON represents a record.
