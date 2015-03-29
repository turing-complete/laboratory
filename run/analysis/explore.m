function explore()
  filename = locate('explore');

  points = h5read(filename, '/points');
  values = h5read(filename, '/values');
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

  if x < 0 || y < 0
    error('Cannot find the sweep dimensions.')
  end

  X = reshape(points(x, :), nn, nn);
  Y = reshape(points(y, :), nn, nn);

  MN = min(values(:));
  MX = max(values(:));

  for z = 1:no
    Z = reshape(values(z, :), nn, nn);

    mn = min(Z(:));
    mx = max(Z(:));

    figure;
    surf(X, Y, Z);
    zlim([MN, MX]);
    title(sprintf('Output %d(Input %d, Input %d), Range %f', z-1, x-1, y-1, mx-mn));
  end
end
