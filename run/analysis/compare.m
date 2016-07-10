function compare(extended, printing)
  use('Interaction');
  if nargin < 1; extended = false; end
  if nargin < 2; printing = false; end
  set(0, 'DefaultTextInterpreter', 'none');
  [files, names] = locate('compare');
  for i = 1:length(files)
    process(files{i}, names{i}, extended, printing);
  end
end

function process(file, name, extended, printing)
  active = h5read(file, '/active');
  oerror = h5read(file, '/observe');
  perror = h5read(file, '/predict');

  nm = size(oerror, 1);
  nk = size(oerror, 2);
  nq = size(oerror, 3);

  from = 1;

  for i = 1:nq
    o = oerror(:, :, i);
    p = perror(:, :, i);

    for j = 1:nm
      t = active(from:end);
      Plot.figure(3 * 204 + 20, 3 * 132 - 10);
      semilogy(t, transpose([o(j, from:end); p(j, from:end)]), ...
        'LineWidth', 2, ...
        'Marker', 'o', ...
        'MarkerSize', 14, ...
        'MarkerFaceColor', 'auto'...
      );
      if length(t) > 1
        Plot.limit([0; t]);
      end
      if printing
        set(gca, ...
          'FontName', 'Times New Roman', ...
          'FontSize', 30, ...
          'YMinorTick', 'off', ...
          'YMinorGrid', 'off');
        set(gcf, ...
          'PaperType', 'A4', ...
          'PaperOrientation', 'landscape', ...
          'PaperPositionMode', 'auto');
        print(sprintf('%s_%d_%d', name, i, j), ...
          '-painters', ...
          '-dpdf', ...
          '-r400');
      else
        Plot.title(sprintf('Case %s, Quantity %d, Metric %d', name, i, j));
        Plot.label('Evaluations', 'log(Error)');
        Plot.legend('Observe', 'Predict');
      end
    end
  end

  if ~extended; return; end

  use('Statistics');

  file = locate('reference', name);
  rvalues = h5read(file{1}, '/values');
  rvalues = rvalues(1:2:end, :);

  file = locate('observe', name);
  ovalues = h5read(file{1}, '/values');
  ovalues = ovalues(1:2:end, :);

  file = locate('predict', name);
  pvalues = h5read(file{1}, '/values');
  pvalues = pvalues(1:2:end, :);

  ns = size(pvalues, 2) / nk;

  pvalues = pvalues(:, (end-ns+1):end);
  ovalues = ovalues(:, 1:active(end));

  for i = 1:nq
    plotDistributions({ ...
      'Reference', rvalues(i, :); ...
      'Observe', ovalues(i, :); ...
      'Predict', pvalues(i, :)}, printing);
    plotDensities({ ...
      'Reference', rvalues(i, :); ...
      'Observe', ovalues(i, :); ...
      'Predict', pvalues(i, :)}, printing);
  end
end

function plotDistributions(sets, ~)
  bins = 100;

  names = sets(:, 1);
  data = sets(:, 2);
  count = length(names);

  Plot.figure(count * 400, 400);
  Plot.title('Histogram');
  for i = 1:count
    subplot(1, count, i);
    hist(data{i}, bins);
    Plot.title(names{i});
  end

  Plot.figure(800, 400);
  Plot.title('CDF');
  hold on;
  for i = 1:count
    ecdf(data{i});
  end
  hold off;
  Plot.legend(names{:});
end

function plotDensities(sets, printing)
  names = sets(:, 1);
  data = sets(:, 2);
  count = length(names);

  Plot.figure(800, 400);
  X = [];
  Y = [];
  hold on;
  for i = 1:count
    ksdensity(data{i});
    line = get(get(gcf, 'Children'), 'Children');
    X = [X, line.XData];
    Y = [Y, line.YData];
  end
  hold off;
  Plot.limit(X, 1.05 * Y);
  if printing
    set(gca, 'FontName', 'Times New Roman', 'FontSize', 30);
  else
    Plot.title('PDF');
  end
  Plot.legend(names{:});
end
