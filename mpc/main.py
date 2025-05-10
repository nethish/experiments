import sys
from mpyc.runtime import mpc

async def main():
    await mpc.start()

    secint = mpc.SecInt()

    # Each party has a salary
    inputs = [60000, 80000, 70000]

    if mpc.pid < len(inputs):
        my_input = secint(inputs[mpc.pid])
    else:
        my_input = None

    print(my_input)
    # All parties participate in the input
    # gather() makes it safe across all parties
    secret_salary = await mpc.gather(mpc.input(my_input, senders=range(len(inputs))))

    # secret_salary is a list with all secret inputs
    total_salary = mpc.sum(secret_salary)
    num_employees = len(inputs)
    average_salary = total_salary / num_employees

    avg = await mpc.output(average_salary)
    print(f"[Party {mpc.pid}] Computed Average Salary: {avg}")

    await mpc.shutdown()

mpc.run(main())
