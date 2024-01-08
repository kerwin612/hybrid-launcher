function change2cn() {
    let xhr = new XMLHttpRequest();
    xhr.onload = function() {
        alert(xhr.status);
        document.title='你好，Hybrid Launcher';
        document.body.innerHTML='系统托盘也切换为中文了';
    };
    xhr.open('GET', '/2cn', true);
    xhr.send();
}
