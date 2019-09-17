<img src="https://github.com/alexanderstephan/gobra/blob/master/assets/gobra.svg.png" width="250" height="77.5" />

### Build instructions: 

---

1. *Dependencies*: e.g on Arch Linux: ``sudo pacman -S ncurses go``

2. *Installation*:

+ Install from repository: ``go get github.com/alexanderstephan/gobra.git``
+ Execute binary: ``"$GOPATH/bin/gobra"``

### How to play

---

Move the snake by using the WASD keys. The more food you collect, the larger the snake grows. If the snakes collides with itself or touches the wall the round is over. Score increases depending on the time it takes you to collect the next food item. 

### Options

---

``-v`` enables `vim` keybindings

``-d`` outputs useful debug information

``-n`` open boundaries
