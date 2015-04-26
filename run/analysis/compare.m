function compare()
  use('Interaction');

  filename = locate('reference');
  rvalues = h5read(filename, '/values');
  rvalues = rvalues(1:2:end, :);

  filename = locate('predict');
  pvalues = h5read(filename, '/values');
  pvalues = pvalues(1:2:end, :);

  filename = locate('compare');
  steps = h5read(filename, '/steps');
  oerror = h5read(filename, '/observe');
  perror = h5read(filename, '/predict');

  nk = size(oerror, 2);
  nq = size(oerror, 3);
  ns = size(pvalues, 2) / nk;

  pvalues = pvalues(:, (end-ns+1):end);

  count = cumsum(steps);

  for i = 1:nq
    o = oerror(:, 2:end, i);
    p = perror(:, 2:end, i);

    Plot.figure(800, 400);
    line(count(2:end), log10(o)');
    line(count(2:end), log10(p)', 'LineStyle', '--');
    legend('Expectation', 'Variance', 'Density');
  end

  for i = 1:nq
    r = rvalues(i, :);
    p = pvalues(i, :);

    [~, ~, error] = kstest2(r, p);

    Plot.figure(800, 400);
    title(sprintf('Histogram (error %.4e)', error));
    subplot(1, 2, 1);
    hist(r, 100);
    title('Reference');
    subplot(1, 2, 2);
    hist(p, 100);
    title('Predict');

    Plot.figure(800, 400);
    title(sprintf('Empirical CDF (error %.4e)', error));
    hold on;
    ecdf(r);
    ecdf(p);
    hold off;
    legend('Reference', 'Predict');
  end
end
