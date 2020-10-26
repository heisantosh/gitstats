# gitstat
A CLI wrapper around git commands to get contributors stats of a repo for all branches.

## Usage
Supported output formats are table/JSON/CSV. Default output format is table.
Stats can be sorted by commits, additions, deletions.

### Example
Running for a local copy of https://github.com/nshepperd/gpt-2.git.
```bash
$ git remote -v
origin	https://github.com/nshepperd/gpt-2.git (fetch)
origin	https://github.com/nshepperd/gpt-2.git (push)
$ gitstats -output csv
Contributor,Commits,Additions,Deletions,Files
Neil Shepperd <nshepperd@google.com>,40,2092,560,19
Jeff Wu <wuthefwasthat@gmail.com>,25,137298,137136,24
Jeff Wu <WuTheFWasThat@gmail.com>,4,478,5,8
Timothy Liu <tlkh@live.co.uk>,3,274,2,2
Ignacio Lopez-Francos <iglopezfrancos@ignacios-mbp.attlocal.net>,2,5,5,2
James B. Pollack <jamesbradenpollack@gmail.com>,2,2,2,2
Nathan Murthy <1788878+natemurthy@users.noreply.github.com>,2,8,2,1
Neil Shepperd <nshepperd@gmail.com>,2,141,0,3
子兎音 <funaox@gmail.com>,2,2,2,2
Anders <oracleliaprojekt@gmail.com>,1,26,0,2
Armaan Bhullar <ArmaanBhullar@users.noreply.github.com>,1,50,1,3
Biranchi <191425+biranchi2018@users.noreply.github.com>,1,3,0,1
Madison May <madison@indico.io>,1,50,0,3
Mathieu Rene <mathieu@zia.ai>,1,6,2,2
Max Woolf <max@minimaxir.com>,1,1,1,1
Memo Akten <memo@users.noreply.github.com>,1,17,9,3
N Shepperd <nshepperd@google.com>,1,0,0,0
Santosh Heigrujam <santosh.hei@gmail.com>,1,123,0,2
Svilen Todorov <sviltodorov@gmail.com>,1,1,0,1
stephan orlowsky <stephan.orlowsky@aperto.com>,1,6,0,1
$
$ gitstats -o table
+------------------------------------------------+---------+-----------+-----------+-------+
|                  CONTRIBUTOR                   | COMMITS | ADDITIONS | DELETIONS | FILES |
+------------------------------------------------+---------+-----------+-----------+-------+
| Neil Shepperd                                  |      40 |      2092 |       560 |    19 |
| <nshepperd@google.com>                         |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| Jeff Wu                                        |      25 |    137298 |    137136 |    24 |
| <wuthefwasthat@gmail.com>                      |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| Jeff Wu                                        |       4 |       478 |         5 |     8 |
| <WuTheFWasThat@gmail.com>                      |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| Timothy Liu <tlkh@live.co.uk>                  |       3 |       274 |         2 |     2 |
+------------------------------------------------+---------+-----------+-----------+-------+
| Ignacio Lopez-Francos                          |       2 |         5 |         5 |     2 |
| <iglopezfrancos@ignacios-mbp.attlocal.net>     |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| James B. Pollack                               |       2 |         2 |         2 |     2 |
| <jamesbradenpollack@gmail.com>                 |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| Nathan Murthy                                  |       2 |         8 |         2 |     1 |
| <1788878+natemurthy@users.noreply.github.com>  |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| Neil Shepperd                                  |       2 |       141 |         0 |     3 |
| <nshepperd@gmail.com>                          |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| 子兎音 <funaox@gmail.com>                      |       2 |         2 |         2 |     2 |
+------------------------------------------------+---------+-----------+-----------+-------+
| Anders                                         |       1 |        26 |         0 |     2 |
| <oracleliaprojekt@gmail.com>                   |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| Armaan Bhullar                                 |       1 |        50 |         1 |     3 |
| <ArmaanBhullar@users.noreply.github.com>       |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| Biranchi                                       |       1 |         3 |         0 |     1 |
| <191425+biranchi2018@users.noreply.github.com> |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| Madison May                                    |       1 |        50 |         0 |     3 |
| <madison@indico.io>                            |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| Mathieu Rene <mathieu@zia.ai>                  |       1 |         6 |         2 |     2 |
+------------------------------------------------+---------+-----------+-----------+-------+
| Max Woolf <max@minimaxir.com>                  |       1 |         1 |         1 |     1 |
+------------------------------------------------+---------+-----------+-----------+-------+
| Memo Akten                                     |       1 |        17 |         9 |     3 |
| <memo@users.noreply.github.com>                |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| N Shepperd                                     |       1 |         0 |         0 |     0 |
| <nshepperd@google.com>                         |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| Santosh Heigrujam                              |       1 |       123 |         0 |     2 |
| <santosh.hei@gmail.com>                        |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| Svilen Todorov                                 |       1 |         1 |         0 |     1 |
| <sviltodorov@gmail.com>                        |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
| stephan orlowsky                               |       1 |         6 |         0 |     1 |
| <stephan.orlowsky@aperto.com>                  |         |           |           |       |
+------------------------------------------------+---------+-----------+-----------+-------+
```

To sort by commits, then by files
```bash
$ gitstats -sort-by=commits
```

## Installation

### Using Go
```bash
go get -v -u github.com/heisantosh/gitstats
```