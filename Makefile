default: all
all: code.tgz a.out.txt b.out.txt c.out.txt d.out.txt e.out.txt f.out.txt
copy: all
	cp code.tgz a.out.txt b.out.txt c.out.txt d.out.txt e.out.txt f.out.txt /Volumes/KIM-Sharing/hashcode/

code.tgz: main.go
	tar cz main.go >$@

doit: main.go
	go build -o $@ main.go

a.out.txt: doit data/a.txt
	time ./doit data/a.txt >$@
b.out.txt: doit data/b.txt
	time ./doit data/b.txt >$@
c.out.txt: doit data/c.txt
	time ./doit data/c.txt >$@
d.out.txt: doit data/d.txt
	time ./doit data/d.txt >$@
e.out.txt: doit data/e.txt
	time ./doit data/e.txt >$@
f.out.txt: doit data/f.txt
	time ./doit data/f.txt >$@
clean:
	rm -rf doit
	rm code.tgz
