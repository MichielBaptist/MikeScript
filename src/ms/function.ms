
function >> test -> int {
    int x;
    0 -> x;
}

function int x >> f -> float {
    float y;
    x -> y;
    y + 1.0 -> y;
}

function string s, float x >> g -> bool {
    string t;
}