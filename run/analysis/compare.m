function compare()
  filename = locate('observe');
  observed = h5read(filename, '/values');
  observed = observed(1:2:end, :);

  filename = locate('predict');
  predicted = h5read(filename, '/values');
  predicted = predicted(1:2:end, :);

  no = size(observed, 1);

  for i = 1:no
    print(observed(i, :), predicted(i, :));
    figure;
    draw(observed(i, :), predicted(i, :));
  end
end
