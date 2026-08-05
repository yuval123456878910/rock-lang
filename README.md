# Rocks proggraming launguage

# Print hello world:
```
print("hello word")
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
func name(int num) (string) {
    ...
}
```
condition: [type] [name]

# Reach
```
reach "hello.ro"
```
reach [path]

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
There is no break. USE `return`

# Build in list functions
```
house NewList = append(<list>, <data>)
```
returns the added data to the list \
```
house NewItem = look(<list>, <pos>)
```
Currentry you cant select the item in the list using [pos]. Insted you do like what i done.
