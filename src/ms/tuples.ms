1, 2, 3 => t;
t[0] >>= print;
t[1] >>= print;
t[2] >>= print;

0 -> t[1];
t[1] >>= print;

(1, 2), (3, 4) => t2;
5 -> t2[0][1];
t2[0] >>= print;
t2[0][0] >>= print;

1, "hello", 3.14, ("world", 42) => t3;
t3 >>= print;
t3[0] >>= print;
t3[1] >>= print;
t3[2] >>= print;
t3[3] >>= print;

[](int, int){(1, 2), (3, 4)} => my_array;

0 => i;
while (i < 2) {
    my_array[i] >>= print;
    my_array[i][0] >>= print;
    my_array[i][1] >>= print;
    i + 1 -> i;
}

type (int, int) point;

var point p1;
p1 >>= print;