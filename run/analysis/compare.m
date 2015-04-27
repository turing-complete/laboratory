function compare(extended)
  use('Interaction');

  filename = locate('compare');
  steps = h5read(filename, '/steps');
  oerror = h5read(filename, '/observe');
  perror = h5read(filename, '/predict');

  nm = size(oerror, 1);
  nk = size(oerror, 2);
  nq = size(oerror, 3);

  count = cumsum(steps);

  labels = {'Expectation', 'Variance', 'Distribution'};
  for i = 1:nq
    o = oerror(:, 2:end, i);
    p = perror(:, 2:end, i);

    Plot.figure(1200, 300);
    for j = 1:nm
      subplot(1, nm, j);
      semilogy(count(2:end), [o(j, :); p(j, :)], 'Marker', 'o');
      legend('Observe', 'Predict');
      Plot.title(labels{j});
    end
  end

  if ~exist('extended', 'var'); return; end

  use('Statistics');

  filename = locate('reference');
  rvalues = h5read(filename, '/values');
  rvalues = rvalues(1:2:end, :);

  filename = locate('observe');
  ovalues = h5read(filename, '/values');
  ovalues = ovalues(1:2:end, :);

  filename = locate('predict');
  pvalues = h5read(filename, '/values');
  pvalues = pvalues(1:2:end, :);

  ns = size(pvalues, 2) / nk;

  pvalues = pvalues(:, (end-ns+1):end);
  ovalues = ovalues(:, 1:count(end));

  for i = 1:nq
    plotDistributions('Reference', rvalues(i, :), 'Observe', ovalues(i, :));
    plotDistributions('Reference', rvalues(i, :), 'Predict', pvalues(i, :));
  end
end

function plotDistributions(name1, data1, name2, data2)
  bins = 100;

  [F1, F2] = distribute(data1, data2);
  error = Error.computeNRMSE(F1, F2);

  Plot.figure(800, 400);
  title(sprintf('Histogram (samples %d, error %.4e)', length(data2), error));
  subplot(1, 2, 1);
  hist(data1, bins);
  title(name1);
  subplot(1, 2, 2);
  hist(data2, bins);
  title(name2);

  Plot.figure(800, 400);
  title(sprintf('Empirical CDF (samples %d, error %.4e)', length(data2), error));
  hold on;
  ecdf(data1);
  ecdf(data2);
  hold off;
  legend(name1, name2);
end

function [cdf1, cdf2] = distribute(data1, data2)
  bins = 100;

  edges = linspace(min(min(data1), min(data2)), max(max(data1), max(data2)), bins + 1);
  edges(end) = Inf;

  cdf1 = histc(data1, edges);
  cdf1 = cdf1(1:end-1);
  cdf1 = cumsum(cdf1) / sum(cdf1);

  cdf2 = histc(data2, edges);
  cdf2 = cdf2(1:end-1);
  cdf2 = cumsum(cdf2) / sum(cdf2);
end
