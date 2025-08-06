
20 => size;

// Init list
[size]int{} => my_list;

// assign values
0 => i;
while i < size {
    i * i -> my_list[i];
    i + 1 -> i;
}

// print values
0 -> i;
while i < size {
    i, "*", i, "=", my_list[i] >>= print;
    i + 1 -> i;
}
