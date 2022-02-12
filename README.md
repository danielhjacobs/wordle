Run `go build`, then run the executable. Requires gcc for the sqlite3 module. 

linux_word_list.txt is a copy of /usr/share/dict/words from a Linux machine

wordle-answers-alphabetical.txt is a copy of https://gist.github.com/cfreshman/a03ef2cba789d8cf00c08f767e0fad7b

wordle-allowed-guesses.txt is a copy of https://gist.github.com/cfreshman/cdcdf777450c5b5301e439061d29694c

If words.db does not exist, it is created from these three files; otherwise, these three files are not needed