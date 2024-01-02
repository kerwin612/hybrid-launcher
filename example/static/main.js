function exit() {
    let xhr = new XMLHttpRequest();
    xhr.onload = function() {
        alert(xhr.status);
    };
    xhr.open('GET', '/exit', true);
    xhr.send();
}
