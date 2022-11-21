let data = []

function addData(event) {
  event.preventDefault()

  let name = document.getElementById("pname").value;
  let sdate = document.getElementById("sdate").value;
  let edate = document.getElementById("edate").value;
  let desc = document.getElementById("desc").value;
  let tech = document.querySelectorAll('input[name="tech"]:checked');
  let gambar = document.getElementById("image").files;

  let techno = []
  tech.forEach((checkbox) => {
    techno.push(checkbox.value)
  })

  if (name == "") {
    return alert('Isi nama proyek!')
  } else if (sdate == "") {
    return alert('Kapan proyek dimulai?')
  } else if (edate == "") {
    return alert('Kapan proyek berakhir?')
  } else if (desc == "") {
    return alert('Isi deskripsi proyek!')
  } else if (gambar.length == 0) {
    return alert('gambar tidak boleh kosong!')
  }

  let image = URL.createObjectURL(gambar[0])

  let project = {
    name,
    duration: lamaWaktu(),
    desc,
    image,
    techno
  }

  data.push(project)
  console.log(data)
  projectPeek()
}

function projectPeek() {
  document.getElementById("projects").innerHTML = ``
  for (let index = 0; index < data.length; index++) {
    document.getElementById("projects").innerHTML +=
    `<div class="prjcard">
    <div class="prjimage">
        <img src="${data[index].image}" style="width: 300px;">
    </div>
        <a href="article.html" target="_blank" class="prjtitle">${data[index].name}</a>
    <div>
        <p class="prjtext">durasi: ${lamaWaktu(data[index].duration)}</p>
    </div>
    <div>
        <p class="prjtext">${data[index].desc}</p>
    </div>
    <div>
        ${(data[index].techno).join(' ')}
    </div>
    <div>
        <button class="button">edit</button>
        <button class="button">delete</button>
    </div>
    </div>`
  }
}

function lamaWaktu(durasi) {
  let sdate = document.getElementById("sdate").value;
  let edate = document.getElementById("edate").value;

  selisih = (Date.parse(edate)) - (Date.parse(sdate))
  harian = Math.floor(selisih / (1000 * 60 * 60 * 24))
  bulanan = Math.floor(selisih / (1000 * 60 * 60 * 24 * 30))
  tahunan = Math.floor(selisih / (1000 * 60 * 60 * 24 * 30 * 12))
  sisahari = harian - (bulanan * 30)
  sisabulan = bulanan - (tahunan * 12)

  if (selisih < 0) {
    return alert('Masukkan tanggal dengan benar!');
  }

  if (harian <= 30) {
    return `${harian} hari`
  } else if (harian <= 365) {
    return `${bulanan} bulan ${sisahari} hari`
  } else {
    return `${tahunan} tahun ${sisabulan} bulan`
  }

  // CODINGAN LAMA
  // let eyear = new Date(edate).getFullYear()
  // let emonth = new Date(edate).getMonth()
  // let eday = new Date(edate).getDate()
  // let syear = new Date(sdate).getFullYear()
  // let smonth = new Date(sdate).getMonth()
  // let sday = new Date(sdate).getDate()

  // let yeardiff = (eyear - syear)
  // let monthdiff = (emonth - smonth)
  // let daydiff = (eday - sday)
  
  // FAILURE: IF looping jika form tidak reset
  // if (yeardiff < 0) {
  //   return alert('Masukkan tahun dengan benar!')
  // } else if ((monthdiff < 0) && (yeardiff <= 0)) {
  //   return alert('Masukkan bulan dengan benar')
  // } else if ((daydiff < 0) && (monthdiff <= 0) && (yeardiff < 0)) {
  //   return alert('Masukkan tanggal dengan benar!')
  // }

  // if (yeardiff > 0) {
  //   return `${yeardiff} tahun`
  // } else if (monthdiff > 0) {
  //   return `${monthdiff} bulan`
  // } else {
  //   return `${daydiff} hari`
  // }
}