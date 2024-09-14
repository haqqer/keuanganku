import { authMe, logout } from "./auth.js";
import { getAll, postTx, getChart, delTx } from "./keuangankuApi.js";

const logoutBtn = document.getElementById("logout");
const addModal = document.getElementById("modal");
const addButton = document.getElementById("addButton");
const tableContentBody = document.getElementById("table-content-body");
const txType = document.getElementById("tx-type");
const txMonth = document.getElementById("tx-month");

const createBtn = document.getElementById("createBtn");

const ctx = document.getElementById("revenueChart").getContext("2d");

let pieChart = new Chart();

const initApp = async () => {
  addButton.addEventListener("click", showAddModal);
  const resAuthMe = await authMe();
  const userinfo = await resAuthMe.json();

  const month = parseInt(new Date().getMonth()) + 1;

  txMonth.value = month;

  await headerUser(userinfo);
  await viewsData(txMonth.value);
  await showChart(txType.value, txMonth.value);
};

var dynamicColors = () => {
  var r = Math.floor(Math.random() * 255);
  var g = Math.floor(Math.random() * 255);
  var b = Math.floor(Math.random() * 255);
  return "rgb(" + r + "," + g + "," + b + ")";
};

const showChart = async (txType, month) => {
  const dataChart = await getChart(txType, month);

  let labels = [];
  let totals = [];
  let colors = [];
  if (dataChart != null) {
    for (var e of dataChart) {
      labels.push(e.category);
      totals.push(e.total);
      colors.push(dynamicColors());
    }
  }

  let typeDisplay = "Expense";
  if (txType == "in") {
    typeDisplay = "Income";
  }

  pieChart = new Chart(ctx, {
    type: "pie",
    data: {
      labels: labels,
      datasets: [
        {
          label: "Data",
          data: totals,
          backgroundColor: colors,
          tension: 0.1,
        },
      ],
    },
    options: {
      responsive: true,
      plugins: {
        legend: {
          position: "bottom",
        },
        title: {
          display: true,
          text: typeDisplay,
        },
      },
    },
  });
};

const headerUser = (userinfo) => {
  const username = document.getElementById("username");
  username.innerText = userinfo.data.name;
  const profileImgUrl = document.getElementById("profileImgUrl");
  profileImgUrl.setAttribute("src", userinfo.data.picture_url);
};

const showAddModal = () => {
  addModal.style.display = "block";
};

window.onclick = function (event) {
  if (event.target == addModal) {
    addModal.style.display = "none";
  }
};

txType.addEventListener("change", async (e) => {
  pieChart.destroy();
  await showChart(txType.value, txMonth.value);
});
txMonth.addEventListener("change", async (e) => {
  pieChart.destroy();
  await showChart(txType.value, txMonth.value);
  await viewsData(txMonth.value);
});

createBtn.addEventListener("click", async (e) => {
  const type = document.getElementById("type-form");
  const category = document.querySelector("input[name='category']");
  const datepickerEl = document.getElementById("default-datepicker");
  const datepicker = new Datepicker(datepickerEl);
  datepicker.setDate(moment().format("MM/DD/YYYY"));
  const currentDate = datepicker.getDate();
  const amount = document.querySelector("input[name='amount']");
  const desc = document.querySelector("input[name='desc']");

  console.log(moment(currentDate));
  const payload = {
    type: type.value.toLowerCase(),
    category: category.value.toLowerCase(),
    date: moment(currentDate).add(1, "days").format(),
    amount: parseInt(amount.value),
    desc: desc.value,
  };
  const result = await postTx(payload);
  category.value = "";
  amount.value = 0;
  desc.value = "";
  setTimeout(async () => {
    pieChart.destroy();
    await initApp();
  }, 500);
});

const delData = async (id) => {
  try {
    await delTx(id);
    pieChart.destroy();
    await initApp();
  } catch (error) {
    console.log(error);
  }
};

const viewsData = async (month) => {
  const txs = await getAll(month);
  let rows = [];
  for (let i = 0; i < txs.length; i++) {
    let delButton = `<td><a href="#" class="delButton px-4 py-2 bg-red-500 rounded-full text-white font-bold hover:bg-red-700" data-id=${txs[i].id}>Delete</a></td>`;
    let amountEl = `<td>${txs[i].amount}</td>`;
    let bgcolor = "";
    if (txs[i].type == "out") {
      bgcolor = "bg-yellow-200";
    } else {
      bgcolor = "bg-emerald-200";
    }

    let row = `<tr class="${bgcolor} font-medium border-b"><td class="px-6 py-4">${txs[i].category}</td><td>${txs[i].desc}</td>${amountEl}${delButton}</tr>`;
    rows.push(row);
  }
  tableContentBody.innerHTML = rows.join("");
  const delButtons = document.querySelectorAll(".delButton");
  delButtons.forEach((btn) => {
    btn.addEventListener("click", async (c) => {
      const target = c.target;
      await delData(target.dataset.id);
    });
  });
};

logoutBtn.addEventListener("click", logout);

document.addEventListener("DOMContentLoaded", initApp);
