function submisi() {
    let name = document.getElementById("nama").value
    let email = document.getElementById("surel").value
    let phone = document.getElementById("nomor").value
    let subject = document.getElementById("subjek").value
    let message = document.getElementById("pesan").value
    console.log(name, email, phone, subject, message);

    if (name == "") {
        return alert('nama harus diisi')
    } else if (email == "") {
        return alert('email harus diisi')
    } else if (phone == "") {
        return alert('nomor telepon harus diisi')
    } else if (subject == "") {
        return alert('subject harus diisi')
    } else if (message == "") {
        return alert('message harus diisi')
    }

    let mailing = document.createElement('a')
    mailing.href = `mailto:${email}?subject=${subject}&body=Hello my name is ${name}, ${subject}, ${message}`;
    mailing.click()
}