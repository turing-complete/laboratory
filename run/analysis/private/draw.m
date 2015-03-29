function draw(observed, predicted)
  [~, ~, delta] = kstest2(observed, predicted);

  subplot(3, 2, 1);
  title(sprintf('CDF (delta %.4e)',delta));
  [y, x] = ksdensity(observed, 'function', 'cdf');
  line(x, y, 'LineStyle', '-');
  [y, x] = ksdensity(predicted, 'function', 'cdf');
  line(x, y, 'LineStyle', '--');
  legend('Observed', 'Predicted');

  subplot(3, 2, 2);
  title('PDF');
  [y, x] = ksdensity(observed, 'function', 'pdf');
  line(x, y, 'LineStyle', '-');
  [y, x] = ksdensity(predicted, 'function', 'pdf');
  line(x, y, 'LineStyle', '--');
  legend('Observed', 'Predicted');

  subplot(3, 2, 3);
  hist(observed, 100);
  title('Observed');
  subplot(3, 2, 4);
  hist(predicted, 100);
  title('Predicted');

  subplot(3, 2, 5);
  h = qqplot(observed, predicted);
  x = h(2).XData;
  y = h(2).YData;
  angle = atan((y(2) - y(1))/(x(2) - x(1))) / pi * 180;
  title(sprintf('Q-Q plot (angle %.4f)', angle))
end
