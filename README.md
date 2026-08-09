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

# List
To assine a list you can only use the var keyword, not house because house will assine all the elements to agiven args in house.
```
house list1 = [1,2,3] # bad-error
var list1 list = [1,2,3] # good
```
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
# Scan
To get an input you use scan keyfunc:
```
house Text = scan(<text>)
```

# Pipe oporation
to pass an arg to a function with a lot of readability, use '|>'!
```
house <name> = <data> |> <function>
```
Example:
```
var list1 list = [1,2,3]
house a = list1 |> look(1)
print(a)
```

# Dict
```
var <name> dict = {<item>:<value>}
```
to get an item you do this:
```
at(<dict_var>, <key>)
```
# How to run?
Download the files, install go and run:
```
go run rock.go <filepath>
```
