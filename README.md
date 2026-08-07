# Rocks proggraming launguage
**This project containes a lot of ai generated content**
# Print hello world:
```
print("hello word")
```

# commets
```
# - for one line
#!
    for multipale lines
!#
```

# Set varubles
```
var a int = 4
const b int = 10
house c = 10 
```
avaible types: string, int, float \
**House** is auto type but always mutible

# Func
```
func name(int num float num2) (string int) {
    ...
}
```
condition: [type] [name] /
no need for "," in perameters and return values

# Reach
```
reach "$/hello.ro"
```
reach [path] \
For current directory you put "$" inside the reach. ONLY IN REACH! \
You can also put a link to a file online:
```
reach "http://...."
```
# Return
```
return 1, "2"
```
return [data]

# Conditions
All condtions are integers: 10 > 1 = 1

# If-Elseif-Else
```
if [condition]  {
    ...
} elseif [condition] {
    ...
} else {
    ...
}
```

# While
```
while [condition] {
    ...
}
```
use `break` to stop

# Build in list functions
```
house NewList = append(<list>, <data>)
```
returns the added data to the list \
```
house NewItem = look(<list>, <pos>)
```
Currentry you cant select the item in the list using [pos]. Insted you do like what i done.

# Thread/Await
Decleration of a thread:
```
House mythread = thread [function]
```
To access the thread return values you do:
```
House [values] = await mythread
```
**YOU ONLY CAN USE THE HOUSE KEYWORD!**