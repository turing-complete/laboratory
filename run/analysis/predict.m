function predict()
  files = locate('predict');
  for i = 1:length(files)
    process(files{i}, grid);
  end
end

function process(file)
  values = h5read(file, '/values');
  values = values(1:2:end, :);

  [~, ns] = size(values);

  I = randi(ns, 1, min(10, ns));

  figure;
  plot(values(:, I), 'Marker', 'o');
  title('Predicted samples');
end
