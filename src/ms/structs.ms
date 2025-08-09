
type struct test {
    []int a;
    float b;
    string c;
}

var test s1;

s1 >>= print;
s1.a >>= print;
s1.b >>= print;
s1.c >>= print;

[10]int{} -> s1.a;
1 -> s1.a[0];
s1.a >>= print;