int x;
int y;
int z;
string s;

while x < 5 {
    0 -> y;
    while y < 5 {
        0 -> z;
        while z < 5 {
            x, y, z >> print;
            if z == 3 {
                break;
            }
            z + 1 -> z;
        }
        y + 1 -> y;
    }
    x + 1 -> x;
}

"hello" >> print;
"world" >> print;