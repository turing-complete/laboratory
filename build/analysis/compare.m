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
      if length(t) > 1; Plot.limit([0; t]); end
      if ~printing
        Plot.title(name);
        Plot.label('Evaluations', 'log(Error)');
        Plot.legend('Observe', 'Predict');
      end
      makeLegible;
      if printing, printOut('%s_%d_%d', name, i, j); end
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

function plotDistributions(sets, printing)
  bins = 100;

  names = sets(:, 1);
  data = sets(:, 2);
  count = length(names);

  Plot.figure(600, 400);
  if ~printing; Plot.title('Histogram'); end
  hold on;
  for i = 1:count
    histogram(data{i}, bins, 'Normalization', 'pdf');
  end
  hold off;
  Plot.legend(names{:});
  makeLegible;

  Plot.figure(600, 400);
  if ~printing; Plot.title('CDF'); end
  hold on;
  for i = 1:count
    [f, x] = ecdf(data{i});
    plot(x, f);
  end
  hold off;
  Plot.legend(names{:});
  makeLegible;
end

function plotDensities(sets, printing)
  names = sets(:, 1);
  data = sets(:, 2);
  count = length(names);

  Plot.figure(600, 400);
  if ~printing; Plot.title('PDF'); end
  X = [];
  F = [];
  hold on;
  for i = 1:count
    [f, x] = ksdensity(data{i});
    plot(x, f);
    X = [X, x];
    F = [F, f];
  end
  hold off;
  Plot.limit(X, 1.05 * F);
  Plot.legend(names{:});
  makeLegible;
end

function makeLegible()
  set(gca, 'FontName', 'Times New Roman', 'FontSize', 30);
end

function printOut(varargin)
    set(gca, ...
      'YMinorTick', 'off', ...
      'YMinorGrid', 'off');
    set(gcf, ...
      'PaperType', 'A4', ...
      'PaperOrientation', 'landscape', ...
      'PaperPositionMode', 'auto');
    print(sprintf(varargin{:}), ...
      '-painters', ...
      '-dpdf', ...
      '-r400');
end
