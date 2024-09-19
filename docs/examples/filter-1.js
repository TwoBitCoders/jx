//Simple Javascript for jx that takes, tests/data/dragons.json as input and 
//outputs all the dragons in the 'Metallic' family.
{ 
    var r = []; 
    for(i = 0; i < x.dragons.length; i++) { 
        let dragon = x.dragons[i];
        if(dragon.family == 'Metallic') {
            r[r.length] = dragon; 
        } 
    } 
    return r; 
}

//These filter examples are to show, that with jx you can write whatever style 
//of js you are comforable with. In this first example we're shooting for 
//explicit, and compatible, it's just plain old Javascript that should compile 
//all the way back to the nineties.
//
//There are only two unusual things you need to be aware of; you need the braces
//and a variable 'x' has been defined with data.
//
//What are the braces for?
//
//Our code is the body of an arrow function, see here:
//
//https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Functions/Arrow_functions
//
//()=>[body]
//
//The body of an arrow function can be either an "expression body" or a "block body"
//
//()=>expression        expression body
//
//or 
//
//()=>{statements}      block body
//
//Since we have a return statement, that dictates that we have a block body and 
//that means we need braces.
//
//Where did 'x' come from?
//
//The variable 'x' has been populated with the result of parsing a JSON value 
//that was input to jx. Think of it like 'document' or 'window' in the browser.
//
//Conceptually, all together it looks like this:
//
//(x)=>{statements}
