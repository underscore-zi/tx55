Cheat building

## Dumping Memory

Load game in PCSX2 without any cheats enabled. Go into the Online Mode menu to load the online mode executable file.
Use PCSX2 to save state in any slot
Go to your PCSX2 folder, find the sstates folder
Find your newly created .ps2 file the naming convention is something like where it ends with .{slot number}.p2s
	SLPM-68524 (E09F13B5).10.p2s

Open the p2s in in archive program. Newer versions of PCSX2 use a uncommon compression format (ZStd, method=93) so your normal program may not work. WinZIP and 7-Zip ZS (not the normal 7zip) do work.

Extract the `eeMemory.bin` file.

## Finding Offsets

The basic process is just finding the right offset for the cheat, and then using a template

So load the eeMemory.bin file in ps2dis, using the defaults when asked. Then invoke `Analyzer > Invoke Analyzer` this will take a little bit. 

### DNAS Offset
	
	- Find the string "NewDnasConnect" using `Edit > Find Pattern`
	- Click on the start of this string in the lower window: https://se.ri0.us/2023-05-13-155045969-db9d4.png
	- Hit "Space" to mark it
	- Hit F3 to find a cross-reference to it
	- This takes me to something that looks like: https://se.ri0.us/2023-05-13-155148781-20959.png
	- If you scroll up a little ways, you should see the marker for the start of a function "FNC_xxxxxxxx" the instruction should be something starting with `addiu sp, sp` (https://se.ri0.us/2023-05-13-155447575-82e71.png)
	- Note the address (first column), in this case 0x0035d6b8
	- This is the start of the where the DNAS patch will go

	DNAS_OFFSET=0035d6b8

### HTTPS Offset

	- Find the string "mgs3sweb" using `Edit > Find Pattern`
	- There will probably be multiple matches. You are looking for the one where multiple URLS exist AND have some space between them (00 bytes)
	- Looking somethign like this: https://se.ri0.us/2023-05-13-160326148-d4f99.png
	- and not this: https://se.ri0.us/2023-05-13-160352756-f5d5a.png
	- Note the start of the URL including the https://
	- In this case 0x0130FA73, you might want to refer to the top panel to get the exact address in memory for the start since the lower one is a bit misleading unless you understand endianness

	HTTPS_OFFSET=0130FA73

 ### Stats Patch
	- Go to the start of the file, hit `g` and then enter `00000000` (eight zeros)
	- Start search again with `Edit > Find pattern`
	- Check the "As hex string box"
	- Enter `12 00 00 12` as the thing to search: https://se.ri0.us/2023-05-13-162858548-d6728.png
	- Hit Okay
	- Hit F5, we want the second occurance
	- It should be somewhere around between 0x190000 and 0x19FFFF I'd guess
	- There should only be one hit starting with 0x19

	- This one is going to use Find Pattern again, except we are searching for a hex string
	- Its the second occurance of this pattern that we want the address for
	- Note this address: https://se.ri0.us/2023-05-13-163053235-ac1c1.png

	STATS_OFFSET=00193618

### Crypto Patch

	- Since we have an address inside of the main MGO segment, we can use it to figure out the crypto one
	- Take the STATS_OFFSET and add 0x64D8
	- This method may be a bit fragile but it should be the start of a function (have a FNC_xxxxxxxx label)
	- BUt its worked for me

	CRYPTO_OFFSET=00199af0

### Master Code Address
	
	- I don't actually remember how to make this one
	- But the idea is you find some code tha tis often run, and you give it the address.
	- https://forum.gamehacking.org/forum/video-game-hacking-and-development/school-of-hacking/1645-ps2-master-codes

## Writing the RAW Cheats

Now that we have the offsets, we can write the core cheats using those addresses. 

.\mastercode_finder_cmd.exe -v MGS3_N.elf and it found a reference to memcpy at 00889848

---

`Edit > Jump to Labeled` (Ctrl+G) and search for `"libpad` one entry should be something like:

 - `"libpad: buffer addr is not 64 byte align. %08x\n"`
 - or `"libpad2: buffer addr is not 64 byte align. %08x\n"`

Go to That label, hit Space to mark, and then F3 to find a cross-reference. Scroll down from there until an `addu` operation: https://se.ri0.us/2023-05-13-180508472-40ff5.png

Master code changes the first digit to 9 of the address:

90183150 00622821

### DNAS

This is going to be a three sequential 4-byte writes. So to determine the left half of the cheats we need the addresses to write to.

1. DNAS_OFFSET (`0035d6b8`)
2. DNAS_OFFSET+4 (`0035d6bC`) 
3. DNAS_OFFSET+8 (`0035d6C0`) 

The left half is the value we want to write, so: `24020000`, `03e00008`, `00000000`

```
0035d6b8 24020000
0035d6bc 03e00008
0035d6C0 00000000
```

And the last step, is to change the first letter of each line to a `2` which indicates a 4-byte write:


```
2035d6b8 24020000
2035d6bc 03e00008
2035d6C0 00000000
```

### Disable Crypto

This one is very similar to the DNAS one, but just two writes(`03e00008`, `00000000`). One at CRYPTO_OFFSET the other at CRYPTO_OFFSET+4.

```
20199af0 03e00008
20199af4 00000000
```

### Stats Bypass

This one is only needed if you're playing on the old Java server codebase, it disables the in-game request for player stats.

```
20193618 00000000
```

### HTTPS to HTTP

This one is the longest cheat, but it does depend on the domain you want to replace the konami one with. We'll be putting "https://savemgo.com/us/\x00". You do need a trailing `00` byte which is what the `\x00` represents.

First we need to convert this string to a hex string, and then swap the endianness. To swap endianness in this case you take 4byte chunks and reverse the bytes so, `12 34 56 78` would become `78 56 34 12` note the bytes swapped order but not the digits inside each byte. If your last chunk isn't 4-bytes, you add zeros to pad it out. 

https://gchq.github.io/CyberChef/#recipe=To_Hex('None',4)Swap_endianness('Hex',4,true)&input=aHR0cDovL3NhdmVtZ28uY29tL3VzLw

The CyberChef recipe gives us this long string:

```
70 74 74 68 73 2f 2f 3a 6d 65 76 61 63 2e 6f 67 75 2f 6d 6f 00 00 2f 73
```

Now lets split it up into our 4byte writes

```
70747468 732f2f3a 6d657661 632e6f67 752f6d6f 00002f73
```

And add our addresses to make the cheat complete which starts with our HTTPS_OFFSET and adds 4 each time

```
2130FA73 70747468 
2130FA77 732f2f3a 
2130FA7B 6d657661 
2130FA7F 632e6f67 
2130FA83 752f6d6f 
2130FA87 00002f73
```

### Conditional Cheat

There is one last "cheat". Its a line added to all the other cheats that conditionally turns the cheat on or off depending on the game state. Its based on the HTTPS_OFFSET but its a 16-bit (2byte) check. The Address is the HTTPS_OFFSET+1, and the value is `7474` (which is the hex form of the string "tt"). SO its basically looking if the "tt" is in memory where we expect it. And the first character of this one is a `D`

```
D130FA74 00007474
```

This cheat is necessary to make the others work when expected. If they apply too early you end up with bugs. So this cheat gets added before EVERY line of all the other cheats. SO the DNAS patch:

```
2035d6b8 24020000
2035d6bc 03e00008
2035d6C0 00000000
```

Becomes:

```
D130FA73 00007474
2035d6b8 24020000
D130FA73 00007474
2035d6bc 03e00008
D130FA73 00007474
2035d6C0 00000000
```


