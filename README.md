# Maze

For this challange we will be creating, viewing and solving mazes.


Binary map format:

```
   // 00 01 02 03  04 05 06 07 
      01                       // byte order marker      -- 1 byte       : little endian
                               //    0 for big endian 
                               //    1 for little endian 

      5D 90                    // magic                   -- 2 bytes     : 5D 90
      01                       // version                 -- 1 byte      : 1

      10 10                    // map size width x height -- 2 bytes     : 16,16 

      20 00                    // map data length(16 bit) -- 2 byte      : 32
      FF FF 80 0D  B7 51 B5 65 // map data (MSb first, packed)
      E5 0D DC B1  98 55 A3 D3 //    a bit set to 1 represents a wall
      87 59 B1 4D  94 55 D3 95 //    a bit set to 0 represents a path
      88 11 B7 ED  80 01 FF FF //    note: one can not travel diagonally

      01 01                    // start x,y position      -- 2 bytes     : 1,1

      02                       // items length            -- 1 byte      : 2 items
                               //    item type value:
                               //       00 required goal  -- 2 bytes
                               //       01 optional goal  -- 2 bytes
                               //       02 warp           -- 4 bytes

                               // item 0
      00                       // required goal           -- 2 bytes     : required goal 
      08 08                    //    x,y position on the map             : 8,8

                               // item 1
      02                       // warp                    -- 4 bytes     : warp 
      0C 0A                    //    x,y position on the map             : 12,10
      01 01                    //    x,y warp destination                : 1,1

```

## Challenge 0: Reader 

1. Create a reader for the above map format

## Challenge 1: Visulization

<ol>
<li>Write a program that given a map file, visualize it to the screen.</li>
<li>
validate the map file with the following rules:
<ol>
 <li>Map data length is properly computed</li>
 <li>Ensure there is a start</li>
 <li>Ensure there is at least one required goal</li>
 <li>Ensure no elements overlap.</li>
</ol>
</li>

<li>
Map:
<ul>
<li>X -- Start position</li>
<li>G -- Goal</li>
<li>W -- Warp</li>
</ul>

```
████████████████
█X           █ █
█ ██ ███ █ █   █
█ ██ █ █ ██  █ █
██   █ █    ██ █
██ ███  █ ██   █
█ █   ████ █  ██
█    ███G█ ██  █
█ ██   █ █  ██ █
█  █ █   █ █W█ █
██ █  ███  █ █ █
█   █      █   █
█ ██ ██████ ██ █
█              █
████████████████
```
</li>
</ol>



### Bonus

 * Render the map to an image.

## Challenge 2: Maze Generation

* Modify the above program to generate a random maze.
* Modify the above program to report if the maze is invalid.

## Challenge 3: Maze Validation

* In addition to the above validation, make sure there is a solution (a first don't worry about maps with warps).
* Make sure there is only one solution to a map.
* Be able to deal with maps with Warps.

## Challenge 4: Maze Solver

* Return the solution of the given map as a set of instructions.
* Create an animated gif showing the how to solve the map

