function compare()
  filename = locate('compare');

  steps = h5read(filename, '/steps');
  observe = h5read(filename, '/observe');
  predict = h5read(filename, '/predict');

  nm = size(observe, 1);
  ns = size(observe, 2);
  nq = size(observe, 3);

  count = cumsum(steps);

  for i = 1:nq
    figure;
    line(count(2:end), log10(observe(:, 2:end, i))');
    line(count(2:end), log10(predict(:, 2:end, i))', 'LineStyle', '--');
  end

  legend('Expectation', 'Variance', 'Density');
end
