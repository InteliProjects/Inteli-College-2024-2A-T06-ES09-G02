#!/bin/bash

# Executa os testes Go
echo "Running Go tests..."
go test ./src/tests/... -v

# Executa os testes JMeter
echo "Running JMeter tests..."
jmeter -n -t ./jmeter/test_plan.jmx -l ./jmeter/results/test_results.csv -e -o ./jmeter/results/test_results.html

echo "Tests completed."
