function observe()
  filename = locate('observe');
  values = h5read(filename, '/values');
  values = values(1:2:end, :);

  [no, ns] = size(values);

  I = randi(ns, 1, min(10, ns));

  figure;
  plot(values(:, I), 'Marker', 'o');
  title('Observed samples');
end
