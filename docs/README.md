# Building the docs

```bash
(mkvirtualenv -r requirements.txt docs)
workon docs
make html
google-chrome _build/html/index.html
```


## Local Server

Use npm to install http-server.

*TODO: expand http-install section*

```bash
http-server _build/html
``` 


## Conversion from Markdown

Markdown is easier to write than .rst format and is already used in README.md files. 

**[m2r](https://pypi.org/project/m2r/)** converts .md to .rst. 

### Installation
```bash
 pip install m2r
```
Or,
```bash
python3 -m pip install m2r
```

###Usage
```bash
m2r your_document.md [your_document2.md ...]
```
Then you will find *your_document.rst* in the same directory.

