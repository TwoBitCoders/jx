//It turns out if you have a single expression for your arrow function body then 
//Javascript will do the return for you. If we don't need the return statements
//then we don't need the braces. After removing the returns and braces we're 
//left with this one-liner that will fit a lot nicer on the command line. 
//There's a couple more characters that could be removed if you're so inclined.
x.dragons.filter(dragon=>dragon.family === "Metallic");

