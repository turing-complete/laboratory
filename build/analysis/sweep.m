function sweep(printing)
  use('Interaction');
  if nargin < 1; printing = false; end
  files = locate('sweep');
  for i = 1:length(files)
    process(files{i}, printing);
  end
end

function process(file, printing)
  points = h5read(file, '/points');
  values = h5read(file, '/values');
  values = values(1:2:end, :);

  ni = size(points, 1);
  no = size(values, 1);
  np = size(values, 2);
  nn = sqrt(np);

  x = -1;
  y = -1;

  for i = 1:ni
    if length(unique(points(i, :))) > 1
      x = i;
      break;
    end
  end

  for i = (x+1):ni
    if length(unique(points(i, :))) > 1
      y = i;
      break;
    end
  end

  if x < 0
    error('Cannot find any sweep dimension.')
  elseif y < 0
    figure;
    line(points(x, :), values);
    if ~printing
      Plot.title('Outputs(Input %d)', x-1);
    end
  else
    X = reshape(points(x, :), nn, nn);
    Y = reshape(points(y, :), nn, nn);

    for z = 1:no
      Z = reshape(values(z, :), nn, nn);

      mn = min(Z(:));
      mx = max(Z(:));

      if printing
        Z = (Z - mn) / (mx - mn);
        mn = 0;
        mx = 1;
      end

      if printing
        Plot.figure(400, 300);
      else
        Plot.figure;
      end

      mesh(X, Y, Z);
      zlim([mn, mx]);
      if printing
        set(gca, ...
          'FontName', 'Times New Roman', ...
          'FontSize', 30);
      else
        Plot.title('Output %d(Input %d, Input %d)', z-1, x-1, y-1);
      end
    end
  end
end
