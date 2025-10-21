document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("order-form");
    const input = document.getElementById("order-id");
    const resultDiv = document.getElementById("order-result");

    form.addEventListener("submit", async function (e) {
        e.preventDefault();
        const orderId = input.value.trim();
        if (!orderId) return;
        try {
            const response = await fetch(`/order?id=${orderId}`);
            const order = await response.json();
            resultDiv.textContent = JSON.stringify(order, null, 2);
        } catch {
            resultDiv.textContent = "Ошибка при получении заказа";
        }
    });
});
