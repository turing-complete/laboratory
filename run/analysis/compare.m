function compare()
  use('Interaction');

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
    o = oerror(:, 2:end, i);
    p = perror(:, 2:end, i);

    Plot.figure(800, 400);
    line(count(2:end), log10(o)');
    line(count(2:end), log10(p)', 'LineStyle', '--');
    legend('Expectation', 'Variance', 'Density');
  end

  for i = 1:nq
    o = ovalues(i, :);
    p = pvalues(i, :);

    [~, ~, delta] = kstest2(o, p);

    Plot.figure(800, 400);
    title(sprintf('CDF (delta %.4e)',delta));
    hold on;
    ecdf(o);
    ecdf(p);
    hold off;
    legend('Observe', 'Predict');

    Plot.figure(800, 400);
    subplot(1, 2, 1);
    hist(o, 100);
    title('Observe');
    subplot(1, 2, 2);
    hist(p, 100);
    title('Predict');
  end
end
