function observe()
  use('Interaction');

  filename = locate('observe');
  values = h5read(filename, '/values');
  values = values(1:2:end, :);

  [no, ns] = size(values);

  figure;
  I = randi(ns, 1, min(10, ns));
  plot(values(:, I), 'Marker', 'o');
  Plot.title('Observed samples');

  for i = 1:no
    figure;
    hist(values(i, :), 100);
    Plot.title('Output %d', i);
  end
end
