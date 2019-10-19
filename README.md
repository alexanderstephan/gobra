## gobra 🐍 

*Feature-rich snake clone in Go using ncurses*

### Build instructions: 

---

#### 1. Dependencies: 

+ On Arch Linux: ``sudo pacman -S ncurses go``

#### 2. Installation:

+ Build repository:

```
make
```
+ Install binary:

```
make install
```

+ Remove binary: 

```
make uninstall
```

### How to play

---

Move the snake by using the *WASD* keys. The more food you collect, the larger the snake grows. If the snakes collides with itself or touches the wall the round is over. Score increases depending on the time it takes you to collect the next food item. 

### Options

---

``-v, --vim`` enables `vim` keybindings

``-d`` outputs useful debug information

``-n`` opens boundaries

``-s`` enables sound
