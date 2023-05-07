function showHamburger() {
    let hamburger = document.getElementById("container-hamburger");
    if (hamburger.style.display === "block") {
      hamburger.style.display = "none";
    } else {
      hamburger.style.display = "block";
    }
}
