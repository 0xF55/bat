## Bat


<p align="center">
  <img src="images/Bat.png" alt="batman" width="500"/>
</p>


**bat** is tool & DSL language for wordlist generating

- Why bat ? :
    * Very fast (Golang Based)
    * Easy syntax
    * Allow many patterns
    * User Friendly
    * Flexible (No Specific Pattern)

## Docs

---

### Escape Charachters

**Space:** `$S`

**$:** `$D`

**@:** `$A`

**=:** `$E`


### Special Variables

**$rndn:length** Returns a random number with specific length

**$rndc:length** Returns a random string with specific length

**$rnds:length** Returns a random special chars with specific length

- Example:
  * $rndn:10 returns 2195325923
  * $rndc:5 returns gAWrl // example
  * $rnds:2 returns !@

### Variables

- Declare a variable

```bat
name=bat
age=25
phone=1337
// nested variable
name2=$name$age$phone // bat251337
test=$name$rnds:3
// , to separate between variable and literal value
name3=$name,qwerty
```


### Modifiers

- Modifiers are characters that change the pattern of the variable

- Example:

```bat
name=aura
upper=$+name
lower=$-upper
zigzag=$~name
reversed=$!name
capitaltized=$^name

@*
```

- `@*` write all variables in the same order to output file

**Output**

```bash
cat wordlist.txt
aura
AURA
aura
aUrA
arua
Aura
```

### Lists

- Declare a list

```bat
list=(1,2,3,4,5,6,7,8)
list2=($rndc:10,$rnds:5)
list3=(qwerty,qazxswedc,abc,6,1)
```

### Write Output

Use `@` Control Character to write to ouput file

- Example :

```bat
@$name$age
@$name$rndn:5
```

### Loops

- Range:

```bat
for i=1..10
  @$name,_$i
end
```

**Output**

```pwsh
Get-Content wordlist.txt  
mohamed_1
mohamed_2
mohamed_3
mohamed_4
mohamed_5
mohamed_6
mohamed_7
mohamed_8
mohamed_9
```

- List

```bat
list=(qwerty,secret,password,leet)
name=Yara

for i=$list
  @$name$i
end
```

**Output**

```bash
cat wordlist.txt
Yaraqwerty
Yarasecret
Yarapassword
Yaraleet
```

- Nested Loops

```bat
for i=1..10
  for j=1..10
    @$i$j
  end
end
```

**Also you can use**

```bat
for i=1..10
for j=1..10
@$i$j
end
end
```

- Iterate Over File

```bat
for i=test.txt
    @hello$i
end
```

**Output**

```bash
cat wordlist.txt
hellofromtest1
hellofromtest2
hellofromtest3
hellofromtest4
hellofromtest5
```

## Examples:

```bat
name=test
for i=1970..2026
  @$name$i
  @$+name$i
  @$~name$i
end

```

**Output**

```bat

test1970
TEST1970
tEsT1970
test1971
TEST1971
tEsT1971
test1972
TEST1972
tEsT1972
test1973
TEST1973
tEsT1973
test1974
TEST1974
tEsT1974
test1975
TEST1975
tEsT1975
test1976
TEST1976
tEsT1976
test1977
TEST1977
tEsT1977
..... // till 2025
```

- Example2:

```bat
name=mohamed

for i=1..1000
    @$name,_$rndn:5
end
```

**Output**


```bat
mohamed_80939
mohamed_15960
mohamed_41304
mohamed_62122
mohamed_76753
mohamed_11165
mohamed_29215
mohamed_46849
mohamed_67054
mohamed_43134
mohamed_55809
mohamed_97010
mohamed_16121
.... // 1000 lines
```

- Example3:

```
name=john
father=doe
for i=1..1000
    @$^name$^father$rndc:3
end
```

**Output**

```
JohnDoeFlX
JohnDoeyAm
JohnDoeYTd
JohnDoewKI
JohnDoeTwN
JohnDoeogc
JohnDoelYb
JohnDoeKbo
JohnDoemJG
JohnDoeLrc
JohnDoeIEh
JohnDoeyQI
JohnDoeWpj
JohnDoeBtg
JohnDoemth
JohnDoeTUp
JohnDoeFFP
... // 1000 line
```
- Example4:

```
name=bat
for y=1970..2026
  for m=1..13
    for d=1..31
    @$name$y:$m:d
    end
  end
end
```
**Output**

```
.\bat.exe -i main.bs
Wordlist: wordlist.txt
BatFile:  main.bs
Generated Lines: 20160
```

- Sample of output

```
bat1970:1:1
bat1970:1:2
bat1970:1:3
bat1970:1:4
bat1970:1:5
bat1970:1:6
bat1970:1:7
bat1970:1:8
bat1970:1:9
bat1970:1:10
bat1970:1:11
... till 2025

```


## Rules

- Space in lists or variables not allowed

**Correct**
```
name=bat
list=(1,2,3,4,5,6)
for i=1..5
end
```
**Incorrect**
```
name = bat
list=(1 , 2 , 3, ,4)
list = (1,2,3,4,5,6)
for i = 1..5
    @$i
end
```

- You can't use @,$,= as literal value, instead use Escape Charachters $D,$A,$E
* Warning: if u use it , the program will be in recursive loop

---

## Download & Build

- You can download precompiled binaries from releases page

- Build Source:

```bash
git clone https://github.com/0xF55/bat.git
cd bat
make build
cd bin
```

- Then run the executable

- Good bye bats ^*^