On the PS2 we don't have an easy way to modify the game to point it towards a server we control and get rid of features we don't like (DNAS). So one way to perform minor game patches it through cheat devices like Codebreaker, Gameshark and Action Replay. PCSX2 and the ps2rd homebrew application also support cheat files.

## Cheat Format

For more complete documentaiton of the format check out: https://macrox.gshi.org/The%20Hacking%20Text.htm#ps2_code_types or https://web.archive.org/web/20070902023442/http://hellion00.thegfcc.com/commands.html

But the basic idea for the cheats we use are they have a single letter dedicated to the operation code, and then the next 7 are an address (without the highest 4-bits/letter) and a value. 

For example: `21020304 01020304` has the opcode of `2` which is a 32bit-write. So it would write the value 0x01020304 to the address `0x01020304`. This is also our most used code, simply writing a 32bit value to a location. We use this with sequential locations to overwrite blocks of memory.

The other two code types/opcodes we use are `C` and `D`. These are both "condition" codes, they check if a condition is true/false and then turn on other lines of the cheat. Either a cheat pack will use `C` at the top of each cheat block, or it will use `D` the first, and then every other line. That is because `C` turns on ALL the cheats that follow it, whereas `D` only enables the next line. Especially when entering these cheats on a PS2 cheat device the time saving of using the `C` code was nice, but it was also only supported on codebreaker and coder I think.


## Condition Check

First, why do we even need conditions. It because the game binary initially loads with all of its game code compressed in memory. Then the first thing the program does is inflates that compressed data back. Cheats don't know about this process so if you apply a memory write too early, before the game has been inflated it will be modifying whatever compressed information is there, which is not the data we want to change. So we need to delay the memory editing until the game is fully loaded. That is what the conditional is for, it checks some piece of memory is present that would only be set once the game was fully loaded and decompressed. 

There was no real rhyme or reason behind the condition checks we used in the cheats. I'd just look in memory for something constant to compare with. So especially for the earliest cheats we write (the `C` based ones) those conditions are pretty random. 

The `D` based condition based ones are a bit more explainable. There is a url in memory starting at `0x01322413`, `D` being a 16bit condition uses that address and then looks for the `tt` characters in `http://` to be present there.

The `C` based conditions, while they will turn off/on entire blocks of cheats, it is a 32bit condition check so if you wanted to adopt the `D` to a `C` cheat, you'd need to adjust the value also. In theory `C1322413 68747470` should work, but I have not tried that. It checks the address for `http` instead of looking for `tt` at +1 that address.

With that in mind, I won't cover the conditions when describing the cheats only the actual memory writes, you can insert whatever condition needed for your cheat device.

# Cheats

## DNAS 

This cheat just injects a `return 0;` into the DNAS call so it instantly returns without doing anything.

```mips
24 02 00 00        addiu $v0, $zero, 0       
03 E0 00 08        jr $ra
00 00 00 00        nop
```

These instructions are written 4 bytes at a time (opcode: `2`) starting with the address `0036b1e8`.

```
2036b1e8 24020000
2036b1ec 03e00008
2036b1f0 00000000
```

## HTTPS

This is the cheat that replaces one of Konami's urls with an http one pointing at savemgo.com.

We use `http://` because it wouldn't be too efficent to get a certificate signed today that the PS2 would accept. I believe many of the authority certs are already expired, and it would be costly even if we could get one. However becauase the game does verify the certificate we can't just fake it either or redirect to one.

The 32bit writes start at 0x01322414, which is actually the first `t` of the `https://` string, and then we write `ttp://savemgo.com/us/mgs3/\x00` over the rest of the url. One thing of note is that endianness matters here. So the actual writes are kinda backwards: `:ptt`, `as//`, `gmev`, `oc.o`, `su/m`, `sgm/`, `/3\x00\x00`. This is also why regardless of region with these cheats the web requests always end up to the `/us/` directory on savemgo.com.

```
21322414 3a707474
21322418 61732f2f
2132241c 676d6576
21322420 6f632e6f
21322424 73752f6d
21322428 73676d2f
2132242c 00002f33 
```

## BYPASS

This cheat should not be used when playing on this gameserver.

This cheat just nops a condition out of the original code, that allows the game to end execution using a 0 status code (success) despite not receiving the correct packets yet. We didn't know the proper response to a request for a players ingame stats, but needed to handle such a request in order to join a game. We made the decision just to patch it out with this cheat and hope for the best. As we now understand the proper response it is not necessary and it breaks the in-game display of stats.

In MIPS `00 00 00 00` is a `nop` or `no operation` instruction. So we just write 4 of those over the condition check.

```
20193a60 00000000
```

## CRYPTO

The original game would encrypt the created game information, probably in an attempt to protect the potentially present game password. I'll be honest, I don't remember this cheat was added super early on and I forgot about it until trying to document this. Rather than dealing with figuring out the encryption being used, we just disabled it completely. This time just a return is written to the top of the encryption/decryption routine to disable in both direction (its also expected that the game settings from the server are encrypted)

```mips
03 E0 00 08        jr $ra
00 00 00 00        nop
```

```
20199f38 03e00008
20199f3c 00000000
```
