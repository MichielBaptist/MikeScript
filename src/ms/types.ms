
type (int, int) point;
type []point point_array;

function (point p) >> translate -> point {
    return (p[0] + 1, p[1] + 2);
}

// expects: (1, 2)
(0, 0) 
>>= translate
=> p 
>>= print;

[]point{(1, 2), (3, 4)} => my_points;
my_points[0] >>= print;
my_points[1] >>= print;

(5, 6) -> my_points[0];
(7, 8) -> my_points[1];

my_points[0] >>= print; // expects: (5, 6)
my_points[1] >>= print; // expects: (7, 8)