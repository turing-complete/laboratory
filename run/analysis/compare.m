function compare()
  filename = locate('observe');
  ovalues = h5read(filename, '/values');
  ovalues = ovalues(1:2:end, :);

  ns = size(ovalues, 2);

  filename = locate('predict');
  pvalues = h5read(filename, '/values');
  pvalues = pvalues(1:2:end, :);

  filename = locate('compare');
  steps = h5read(filename, '/steps');
  oerror = h5read(filename, '/observe');
  perror = h5read(filename, '/predict');

  nm = size(oerror, 1);
  nk = size(oerror, 2);
  nq = size(oerror, 3);

  pvalues = pvalues(:, (end-ns+1):end);

  count = cumsum(steps);

  for i = 1:nq
    figure;
    line(count(2:end), log10(oerror(:, 2:end, i))');
    line(count(2:end), log10(perror(:, 2:end, i))', 'LineStyle', '--');
    legend('Expectation', 'Variance', 'Density');

    figure;
    draw(ovalues(i, :), pvalues(i, :));
    print(ovalues(i, :), pvalues(i, :));
  end
end
