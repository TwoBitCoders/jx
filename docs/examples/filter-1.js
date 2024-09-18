//Write a simple Javascript function that takes dragons.json as input and outputs
//all the dragons in the 'Metallic' family
//We make no effort to be concise, just write things out in plain old 
//Javascript.
//Our code is the body of an arrow function, see here:
//https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Functions/Arrow_functions
//()=>[body]
//The body of an arrow function can be either an "expression body" or a "block body"
//()=>expression        expression body
//or 
//()=>{statements}      block body
//Since we have a return statement, that dictates that we have a block body and 
//that means we need braces
//Remember that in jx we are always provided the variable 'x' which is all the JSON data
//that we have to work with
{ 
    var r = []; 
    for(i = 0; i < x.dragons.length; i++) { 
        let dragon = x.dragons[i];
        if(dragon.family === 'Metallic') {
            r.push(dragon); 
        } 
    } 
    return r; 
}

